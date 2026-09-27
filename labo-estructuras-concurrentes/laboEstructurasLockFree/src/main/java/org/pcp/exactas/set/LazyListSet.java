package org.pcp.exactas.set;

import java.util.concurrent.locks.ReentrantLock;

public class LazyListSet implements ConcurrentIntSet {
    private final Node head = new Node(Long.MIN_VALUE);
    private final Node tail = new Node(Long.MAX_VALUE);

    public LazyListSet() { head.next = tail; }

    @Override public boolean add(int value) {
        while (true) {
            Node pred = head, curr = head.next;
            while(curr.key < value){
                pred = curr; 
                curr = curr.next;
            }
            pred.lock.lock();
            try{
                curr.lock.lock();
                try{
                    if (edgeExists(pred, curr)){
                        if(curr.key == value){
                            return false;
                        }
                        else{
                            Node node = new Node(value);
                            node.next = curr; 
                            pred.next = node;
                            return true;
                        }
                    }
                }
                finally{
                    curr.lock.unlock();
                }
            }
            finally{
                pred.lock.unlock();
            }
        }
    }

    @Override public boolean remove(int value) {
        while (true) {
            Node pred = head, curr = head.next;
            while(curr.key < value){
                pred = curr; 
                curr = curr.next;
            }
            pred.lock.lock();
            try{
                curr.lock.lock();
                try{
                    if (edgeExists(pred, curr)){
                        if(curr.key == value){
                            curr.marked = true;
                            pred.next = curr.next;
                            return true;
                        }
                        else{
                            return false;
                        }
                    }
                }
                finally{
                    curr.lock.unlock();
                }
            }
            finally{
                pred.lock.unlock();
            }
        }
    }

    @Override public boolean contains(int value) {
        Node curr = head.next;
        while (curr.key < value) {
            curr = curr.next;
        }
        return curr.key == value && !curr.marked;
    }

    private static boolean edgeExists(Node pred, Node curr) {
        return !pred.marked && !curr.marked && pred.next == curr;
    }

    @Override public boolean isImplemented() {
        return true;
    }

    private static final class Node {
        final long key;
        volatile Node next;
        volatile boolean marked;
        final ReentrantLock lock = new ReentrantLock();
        Node(long key) { this(key, null); } Node(long key, Node next) { this.key=key; this.next=next; }
    }
}
